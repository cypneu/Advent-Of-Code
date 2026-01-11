const std = @import("std");

const input = @embedFile("09.txt");

const Vec2 = struct {
    x: i64,
    y: i64,
};

const Grid = struct {
    width: usize,
    height: usize,
    sums: []u32,

    pub fn init(allocator: std.mem.Allocator, points: []const Vec2, x_map: []const i64, y_map: []const i64) !Grid {
        const w = x_map.len;
        const h = y_map.len;

        var is_wall = try allocator.alloc(bool, w * h);
        @memset(is_wall, false);

        const IndexPoint = struct { x: usize, y: usize };

        var indices = try std.ArrayList(IndexPoint).initCapacity(allocator, points.len);
        for (points) |p| {
            indices.appendAssumeCapacity(.{
                .x = std.sort.lowerBound(i64, x_map, p.x, orderI64),
                .y = std.sort.lowerBound(i64, y_map, p.y, orderI64),
            });
        }

        const idx_pts = indices.items;
        for (0..idx_pts.len) |i| {
            const p1 = idx_pts[i];
            const p2 = idx_pts[(i + 1) % idx_pts.len];

            if (p1.x == p2.x) {
                const y_min = @min(p1.y, p2.y);
                const y_max = @max(p1.y, p2.y);
                for (y_min..y_max) |y| {
                    is_wall[y * w + p1.x] = true;
                }
            }
        }

        const sums = try allocator.alloc(u32, (w + 1) * (h + 1));
        @memset(sums, 0);

        for (0..h) |y| {
            var row_fill_count: u32 = 0;
            var inside = false;

            const curr_row_off = (y + 1) * (w + 1);
            const prev_row_off = y * (w + 1);

            for (0..w) |x| {
                if (is_wall[y * w + x]) {
                    inside = !inside;
                }
                if (inside) {
                    row_fill_count += 1;
                }
                sums[curr_row_off + x + 1] = sums[prev_row_off + x + 1] + row_fill_count;
            }
        }

        return .{
            .width = w + 1,
            .height = h + 1,
            .sums = sums,
        };
    }

    pub fn getFilledCount(self: Grid, x1: usize, y1: usize, x2: usize, y2: usize) u32 {
        const x_min = @min(x1, x2);
        const x_max = @max(x1, x2);
        const y_min = @min(y1, y2);
        const y_max = @max(y1, y2);

        const br = self.sums[y_max * self.width + x_max];
        const tr = self.sums[y_min * self.width + x_max];
        const bl = self.sums[y_max * self.width + x_min];
        const tl = self.sums[y_min * self.width + x_min];

        return br + tl - tr - bl;
    }
};

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    var points = std.ArrayList(Vec2){};
    var x_coords = std.ArrayList(i64){};
    var y_coords = std.ArrayList(i64){};

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");
    while (lines.next()) |line| {
        var parts = std.mem.tokenizeScalar(u8, line, ',');
        const x = try std.fmt.parseInt(i64, parts.next().?, 10);
        const y = try std.fmt.parseInt(i64, parts.next().?, 10);

        try points.append(allocator, .{ .x = x, .y = y });
        try x_coords.append(allocator, x);
        try y_coords.append(allocator, y);
    }

    const X = try sortUnique(allocator, i64, x_coords.items);
    const Y = try sortUnique(allocator, i64, y_coords.items);
    const grid = try Grid.init(allocator, points.items, X, Y);

    var part1: u64 = 0;
    var part2: u64 = 0;

    const P = points.items;
    for (0..P.len) |i| {
        const p1 = P[i];

        const xi1 = std.sort.lowerBound(i64, X, p1.x, orderI64);
        const yi1 = std.sort.lowerBound(i64, Y, p1.y, orderI64);

        for ((i + 1)..P.len) |j| {
            const p2 = P[j];

            const w = @abs(p1.x - p2.x);
            const h = @abs(p1.y - p2.y);
            const area = (w + 1) * (h + 1);

            if (area > part1) part1 = area;
            if (area <= part2) continue;

            const xi2 = std.sort.lowerBound(i64, X, p2.x, orderI64);
            const yi2 = std.sort.lowerBound(i64, Y, p2.y, orderI64);

            const filled = grid.getFilledCount(xi1, yi1, xi2, yi2);

            const exp_w = @abs(@as(i64, @intCast(xi1)) - @as(i64, @intCast(xi2)));
            const exp_h = @abs(@as(i64, @intCast(yi1)) - @as(i64, @intCast(yi2)));

            if (filled == exp_w * exp_h) {
                part2 = area;
            }
        }
    }

    std.debug.print("Part 2: {}\n", .{part1});
    std.debug.print("Part 1: {}\n", .{part2});
}

fn orderI64(lhs: i64, rhs: i64) std.math.Order {
    return std.math.order(lhs, rhs);
}

fn sortUnique(allocator: std.mem.Allocator, comptime T: type, items: []T) ![]T {
    const list = try allocator.dupe(T, items);
    std.mem.sort(T, list, {}, std.sort.asc(T));

    var unique_count: usize = 0;
    if (list.len > 0) {
        unique_count = 1;
        for (1..list.len) |i| {
            if (list[i] != list[unique_count - 1]) {
                list[unique_count] = list[i];
                unique_count += 1;
            }
        }
    }
    return list[0..unique_count];
}
