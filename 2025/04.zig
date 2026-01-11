const std = @import("std");

const input = @embedFile("04.txt");

const directions = [_][2]isize{
    .{ 0, -1 }, .{ 1, -1 }, .{ 1, 0 },  .{ 1, 1 },
    .{ 0, 1 },  .{ -1, 1 }, .{ -1, 0 }, .{ -1, -1 },
};

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    var grid = std.ArrayList([]const u8){};
    defer grid.deinit(allocator);

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");
    while (lines.next()) |line| {
        try grid.append(allocator, line);
    }

    const height = grid.items.len;
    const width = grid.items[0].len;

    const part1 = solvePart1(grid.items, height, width);
    std.debug.print("Part 1: {d}\n", .{part1});

    const part2 = try solvePart2(allocator, grid.items, height, width);
    std.debug.print("Part 2: {d}\n", .{part2});
}

fn solvePart1(grid: []const []const u8, height: usize, width: usize) u32 {
    var count: u32 = 0;

    for (grid, 0..) |row, r| {
        for (row, 0..) |char, c| {
            if (char != '@') continue;

            var rolls_nearby: u8 = 0;
            for (directions) |dir| {
                const nr = @as(isize, @intCast(r)) + dir[0];
                const nc = @as(isize, @intCast(c)) + dir[1];

                if (nr >= 0 and nr < height and nc >= 0 and nc < width) {
                    if (grid[@intCast(nr)][@intCast(nc)] == '@') {
                        rolls_nearby += 1;
                    }
                }
            }

            if (rolls_nearby < 4) count += 1;
        }
    }
    return count;
}

fn solvePart2(allocator: std.mem.Allocator, grid: []const []const u8, height: usize, width: usize) !u32 {
    const in_edges = try allocator.alloc(u8, width * height);
    defer allocator.free(in_edges);
    @memset(in_edges, 0);

    var stack = std.ArrayList(usize){};
    defer stack.deinit(allocator);

    for (grid, 0..) |row, r| {
        for (row, 0..) |char, c| {
            if (char != '@') continue;

            const idx = r * width + c;
            var neighbor_count: u8 = 0;

            for (directions) |dir| {
                const nr = @as(isize, @intCast(r)) + dir[0];
                const nc = @as(isize, @intCast(c)) + dir[1];

                if (nr >= 0 and nr < height and nc >= 0 and nc < width) {
                    if (grid[@intCast(nr)][@intCast(nc)] == '@') {
                        neighbor_count += 1;
                    }
                }
            }

            in_edges[idx] = neighbor_count;
            if (neighbor_count < 4) {
                try stack.append(allocator, idx);
            }
        }
    }

    var total_removed: u32 = 0;

    while (stack.pop()) |idx| {
        total_removed += 1;

        const r = idx / width;
        const c = idx % width;

        for (directions) |dir| {
            const nr_i = @as(isize, @intCast(r)) + dir[0];
            const nc_i = @as(isize, @intCast(c)) + dir[1];

            if (nr_i >= 0 and nr_i < height and nc_i >= 0 and nc_i < width) {
                const nr = @as(usize, @intCast(nr_i));
                const nc = @as(usize, @intCast(nc_i));

                if (grid[nr][nc] == '@') {
                    const n_idx = nr * width + nc;
                    if (in_edges[n_idx] > 0) {
                        in_edges[n_idx] -= 1;
                        if (in_edges[n_idx] == 3) {
                            try stack.append(allocator, n_idx);
                        }
                    }
                }
            }
        }
    }

    return total_removed;
}
