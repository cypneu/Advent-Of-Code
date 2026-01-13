const std = @import("std");

const input = @embedFile("10.txt");

const Machine = struct {
    lights_target: u64,
    num_lights: usize,

    joltage_targets: []const i64,

    buttons: []const Button,

    const Button = struct {
        indices: []const usize,
        light_mask: u64,
    };
};

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    const machines = try parseInput(allocator, input);

    var part1: usize = 0;
    var part2: i64 = 0;

    for (machines) |machine| {
        part1 += try solvePart1(allocator, machine);
        part2 += try solvePart2(allocator, machine);
    }

    std.debug.print("Part 1: {}\n", .{part1});
    std.debug.print("Part 2: {}\n", .{part2});
}

fn parseInput(allocator: std.mem.Allocator, data: []const u8) ![]Machine {
    var machines = std.ArrayList(Machine){};
    var lines = std.mem.tokenizeAny(u8, data, "\r\n");

    while (lines.next()) |line| {
        var lights_target: u64 = 0;
        var num_lights: usize = 0;

        var joltage_targets = std.ArrayList(i64){};
        var button_indices = std.ArrayList([]const usize){};

        var parts = std.mem.tokenizeScalar(u8, line, ' ');
        while (parts.next()) |part| {
            switch (part[0]) {
                '[' => {
                    const lights = try parseLightsDiagram(part);
                    lights_target = lights.target;
                    num_lights = lights.len;
                },
                '(' => {
                    const idxs = try parseIndexList(allocator, part);
                    try button_indices.append(allocator, idxs);
                },
                '{' => try parseTargetsInto(allocator, part, &joltage_targets),
                else => {},
            }
        }

        var buttons = std.ArrayList(Machine.Button){};
        for (button_indices.items) |idxs| {
            const mask = try buildLightMask(idxs);
            try buttons.append(allocator, .{
                .indices = idxs,
                .light_mask = mask,
            });
        }

        try machines.append(allocator, .{
            .lights_target = lights_target,
            .num_lights = num_lights,
            .joltage_targets = try joltage_targets.toOwnedSlice(allocator),
            .buttons = try buttons.toOwnedSlice(allocator),
        });
    }

    return machines.items;
}

const LightsSpec = struct { target: u64, len: usize };

fn parseLightsDiagram(token: []const u8) !LightsSpec {
    const content = token[1 .. token.len - 1];
    var target: u64 = 0;
    for (content, 0..) |ch, i| {
        if (ch == '#') target |= (@as(u64, 1) << @intCast(i));
    }
    return .{ .target = target, .len = content.len };
}

fn parseTargetsInto(
    allocator: std.mem.Allocator,
    token: []const u8,
    out: *std.ArrayList(i64),
) !void {
    const content = token[1 .. token.len - 1];
    var vals = std.mem.tokenizeScalar(u8, content, ',');
    while (vals.next()) |v| {
        try out.append(allocator, try std.fmt.parseInt(i64, v, 10));
    }
}

fn parseIndexList(allocator: std.mem.Allocator, token: []const u8) ![]const usize {
    const content = token[1 .. token.len - 1];

    var result: std.ArrayList(usize) = .{};
    var it = std.mem.tokenizeScalar(u8, content, ',');
    while (it.next()) |raw| {
        try result.append(allocator, try std.fmt.parseInt(usize, raw, 10));
    }

    return try result.toOwnedSlice(allocator);
}

fn buildLightMask(idxs: []const usize) !u64 {
    var mask: u64 = 0;
    for (idxs) |i| {
        mask |= (@as(u64, 1) << @intCast(i));
    }
    return mask;
}

fn collectFreeVariables(allocator: std.mem.Allocator, pivot_row_for_col: []const ?usize) ![]usize {
    var free = std.ArrayList(usize){};
    for (pivot_row_for_col, 0..) |pivot_row, col| {
        if (pivot_row == null) try free.append(allocator, col);
    }
    return try free.toOwnedSlice(allocator);
}

fn solvePart1(allocator: std.mem.Allocator, m: Machine) !usize {
    const rows = m.num_lights;
    const cols = m.buttons.len;

    var matrix = try allocator.alloc(u64, rows);
    var target = try allocator.alloc(u1, rows);

    for (0..rows) |row| {
        target[row] = @intCast((m.lights_target >> @intCast(row)) & 1);

        var row_mask: u64 = 0;
        for (m.buttons, 0..) |btn, col| {
            if (((btn.light_mask >> @intCast(row)) & 1) == 1) {
                row_mask |= @as(u64, 1) << @intCast(col);
            }
        }
        matrix[row] = row_mask;
    }

    const pivot = try rrefGF2(allocator, matrix, target, cols);

    const free_variables = try collectFreeVariables(allocator, pivot);
    const num_free = free_variables.len;

    const combos: usize = @as(usize, 1) << @intCast(num_free);
    var best: usize = std.math.maxInt(usize);
    for (0..combos) |i| {
        var sol_bits: u64 = 0;

        for (free_variables, 0..) |col, bit| {
            if (((i >> @intCast(bit)) & 1) == 1)
                sol_bits |= (@as(u64, 1) << @intCast(col));
        }

        for (0..cols) |col| {
            const pivot_row = pivot[col] orelse continue;
            const parity: u1 = @intCast(@popCount(matrix[pivot_row] & sol_bits) & 1);
            const val = target[pivot_row] ^ parity;
            if (val == 1) sol_bits |= (@as(u64, 1) << @intCast(col));
        }

        const presses: usize = @intCast(@popCount(sol_bits));
        if (presses < best) best = presses;
    }

    return best;
}

fn rrefGF2(allocator: std.mem.Allocator, matrix: []u64, target: []u1, cols: usize) ![]?usize {
    var pivot_row_for_col = try allocator.alloc(?usize, cols);
    @memset(pivot_row_for_col, null);

    const rows = matrix.len;
    var pivot_row: usize = 0;

    for (0..cols) |col| {
        var selected_row: ?usize = null;
        for (pivot_row..rows) |row| {
            if (((matrix[row] >> @intCast(col)) & 1) == 1) {
                selected_row = row;
                break;
            }
        }
        if (selected_row == null) continue;

        std.mem.swap(u64, &matrix[pivot_row], &matrix[selected_row.?]);
        std.mem.swap(u1, &target[pivot_row], &target[selected_row.?]);

        for (0..rows) |row| {
            if (row == pivot_row) continue;
            if (((matrix[row] >> @intCast(col)) & 1) == 1) {
                matrix[row] ^= matrix[pivot_row];
                target[row] ^= target[pivot_row];
            }
        }

        pivot_row_for_col[col] = pivot_row;
        pivot_row += 1;
    }

    return pivot_row_for_col;
}

fn solvePart2(allocator: std.mem.Allocator, machine: Machine) !i64 {
    const rows = machine.joltage_targets.len;
    const cols = machine.buttons.len;

    var matrix_storage = try allocator.alloc(i128, rows * cols);
    @memset(matrix_storage, 0);

    var matrix = try allocator.alloc([]i128, rows);
    for (0..rows) |row| {
        matrix[row] = matrix_storage[row * cols .. (row + 1) * cols];
    }

    var target = try allocator.alloc(i128, rows);
    for (0..rows) |row| target[row] = machine.joltage_targets[row];

    for (machine.buttons, 0..) |button, col| {
        for (button.indices) |row| {
            matrix[row][col] = 1;
        }
    }

    const pivot = try rrefIntegerFractionFree(allocator, matrix, target);
    const free_cols = try collectFreeVariables(allocator, pivot);
    const bounds = try computeFreeVarBounds(allocator, machine, free_cols);

    const solution = try allocator.alloc(i64, cols);
    @memset(solution, 0);

    const context = SolverContext{
        .matrix = matrix,
        .targets = target,
        .pivots = pivot,
        .free_cols = free_cols,
        .bounds = bounds,
        .solution = solution,
    };

    const result = minimizeButtonPresses(context, 0, 0, null);
    return result orelse -1;
}

fn gcdAbsI128(a: i128, b: i128) i128 {
    const ua: u128 = @abs(a);
    const ub: u128 = @abs(b);
    return @as(i128, @intCast(std.math.gcd(ua, ub)));
}

fn rrefIntegerFractionFree(allocator: std.mem.Allocator, matrix: [][]i128, target: []i128) ![]?usize {
    const rows = matrix.len;
    const cols = matrix[0].len;

    var pivot_row_for_col = try allocator.alloc(?usize, cols);
    @memset(pivot_row_for_col, null);

    var pivot_row: usize = 0;
    for (0..cols) |col| {
        var selected_row: ?usize = null;
        for (pivot_row..rows) |row| {
            if (matrix[row][col] != 0) {
                selected_row = row;
                break;
            }
        }
        if (selected_row == null) continue;

        std.mem.swap([]i128, &matrix[pivot_row], &matrix[selected_row.?]);
        std.mem.swap(i128, &target[pivot_row], &target[selected_row.?]);

        const pivot_value = matrix[pivot_row][col];

        for (0..rows) |row| {
            if (row == pivot_row) continue;
            const value = matrix[row][col];
            if (value == 0) continue;

            const g = gcdAbsI128(value, pivot_value);
            const mult_r = @divExact(pivot_value, g);
            const mult_p = @divExact(value, g);

            for (0..cols) |c| {
                matrix[row][c] = matrix[row][c] * mult_r - matrix[pivot_row][c] * mult_p;
            }
            target[row] = target[row] * mult_r - target[pivot_row] * mult_p;
        }

        pivot_row_for_col[col] = pivot_row;
        pivot_row += 1;
    }

    return pivot_row_for_col;
}

fn computeFreeVarBounds(allocator: std.mem.Allocator, machine: Machine, free_cols: []const usize) ![]i64 {
    var bounds = std.ArrayList(i64){};

    for (free_cols) |col| {
        const button = machine.buttons[col];

        var limit: i64 = std.math.maxInt(i64);
        for (button.indices) |row| {
            const target = machine.joltage_targets[row];
            if (target < limit) limit = target;
        }

        try bounds.append(allocator, limit);
    }

    return bounds.toOwnedSlice(allocator);
}

const SolverContext = struct {
    matrix: [][]i128,
    targets: []const i128,
    pivots: []const ?usize,
    free_cols: []const usize,
    bounds: []const i64,
    solution: []i64,
};

fn minimizeButtonPresses(ctx: SolverContext, free_idx: usize, current_sum: i64, limit: ?i64) ?i64 {
    var best = limit;

    if (best != null and current_sum >= best.?) return best;

    if (free_idx == ctx.free_cols.len) {
        if (calculateDependentSum(ctx)) |dep_sum| {
            const total = current_sum + dep_sum;
            if (best == null or total < best.?) {
                return total;
            }
        }
        return best;
    }

    const col = ctx.free_cols[free_idx];
    const max_val = ctx.bounds[free_idx];

    var val: i64 = 0;
    while (val <= max_val) : (val += 1) {
        ctx.solution[col] = val;
        best = minimizeButtonPresses(ctx, free_idx + 1, current_sum + val, best);
    }

    return best;
}

fn calculateDependentSum(ctx: SolverContext) ?i64 {
    var sum: i64 = 0;

    for (ctx.pivots, 0..) |pivot_row_opt, col| {
        const row = pivot_row_opt orelse continue;

        var numerator = ctx.targets[row];
        for (ctx.free_cols) |free_col| {
            numerator -= ctx.matrix[row][free_col] * @as(i128, ctx.solution[free_col]);
        }

        const denominator = ctx.matrix[row][col];
        if (@rem(numerator, denominator) != 0) return null;

        const val_128 = @divExact(numerator, denominator);
        if (val_128 < 0 or val_128 > std.math.maxInt(i64)) return null;

        sum += @intCast(val_128);
    }

    return sum;
}
