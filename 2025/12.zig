const std = @import("std");

const input = @embedFile("12.txt");

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    var shapes = std.ArrayList(Shape){};
    var shape_buffer = std.ArrayList(Point){};
    var current_y: i8 = 0;

    var solved_count: usize = 0;
    var lines = std.mem.tokenizeAny(u8, input, "\r\n");

    while (lines.next()) |line| {
        if (Parser.isQuery(line)) {
            if (shape_buffer.items.len > 0) {
                const s = try Geometry.createShape(allocator, shape_buffer.items);
                try shapes.append(allocator, s);
                shape_buffer.clearRetainingCapacity();
                current_y = 0;
            }

            const query = try Parser.parseQuery(allocator, line, shapes.items);

            if (try Solver.isSolvable(allocator, query, shapes.items)) {
                solved_count += 1;
            }
        } else if (Parser.isHeader(line)) {
            if (shape_buffer.items.len > 0) {
                const s = try Geometry.createShape(allocator, shape_buffer.items);
                try shapes.append(allocator, s);
                shape_buffer.clearRetainingCapacity();
            }
            current_y = 0;
        } else {
            try Parser.parseShapeRow(allocator, &shape_buffer, line, current_y);
            current_y += 1;
        }
    }

    std.debug.print("Part 1: {d}\n", .{solved_count});
}

const Point = struct {
    x: i8,
    y: i8,

    pub fn less(_: void, a: Point, b: Point) bool {
        if (a.y != b.y) return a.y < b.y;
        return a.x < b.x;
    }
};

const Shape = struct {
    variants: []const Variant,
    area: usize,
};

const Variant = struct {
    points: []const Point,
    width: u8,
    height: u8,
};

const Query = struct {
    w: u8,
    h: u8,
    items: []const usize,
    total_area: usize,
};

const Solver = struct {
    pub fn isSolvable(alloc: std.mem.Allocator, q: Query, shapes: []const Shape) !bool {
        const region_area = @as(usize, q.w) * q.h;
        if (q.total_area > region_area) return false;

        const grid = try alloc.alloc(bool, region_area);
        @memset(grid, false);

        const placed = try alloc.alloc(usize, q.items.len);
        @memset(placed, 0);

        return backtrack(grid, q.w, q.h, q.items, 0, shapes, placed);
    }

    fn backtrack(grid: []bool, w: u8, h: u8, items: []const usize, idx: usize, shapes: []const Shape, placed: []usize) bool {
        if (idx == items.len) return true;

        const shape_idx = items[idx];
        const shape = shapes[shape_idx];

        const start_pos = if (idx > 0 and items[idx - 1] == shape_idx) placed[idx - 1] else 0;
        const limit_pos = grid.len;

        for (shape.variants) |variant| {
            if (variant.width > w or variant.height > h) continue;

            const limit_x = w - variant.width;
            const limit_y = h - variant.height;

            var pos = start_pos;
            while (pos < limit_pos) : (pos += 1) {
                const cy = @as(u8, @intCast(pos / w));
                const cx = @as(u8, @intCast(pos % w));

                if (cy > limit_y) break;
                if (cx > limit_x) continue;

                if (canFit(grid, w, cx, cy, variant.points)) {
                    toggle(grid, w, cx, cy, variant.points, true);
                    placed[idx] = pos;

                    if (backtrack(grid, w, h, items, idx + 1, shapes, placed)) return true;

                    toggle(grid, w, cx, cy, variant.points, false);
                }
            }
        }
        return false;
    }

    fn canFit(grid: []const bool, w: u8, cx: u8, cy: u8, points: []const Point) bool {
        for (points) |p| {
            const px = cx + @as(u8, @intCast(p.x));
            const py = cy + @as(u8, @intCast(p.y));
            if (grid[@as(usize, py) * w + px]) return false;
        }
        return true;
    }

    fn toggle(grid: []bool, w: u8, cx: u8, cy: u8, points: []const Point, state: bool) void {
        for (points) |p| {
            const px = cx + @as(u8, @intCast(p.x));
            const py = cy + @as(u8, @intCast(p.y));
            grid[@as(usize, py) * w + px] = state;
        }
    }
};

const Geometry = struct {
    pub fn createShape(alloc: std.mem.Allocator, raw: []const Point) !Shape {
        var variants = std.ArrayList(Variant){};

        for (0..2) |flip_idx| {
            for (0..4) |rot| {
                const flip = (flip_idx == 1);
                const v = try createVariant(alloc, raw, rot, flip);
                if (!isDuplicate(variants.items, v)) {
                    try variants.append(alloc, v);
                }
            }
        }

        return .{ .variants = try variants.toOwnedSlice(alloc), .area = raw.len };
    }

    fn createVariant(alloc: std.mem.Allocator, raw: []const Point, rot: usize, flip: bool) !Variant {
        const points = try alloc.alloc(Point, raw.len);

        var min_x: i8 = 127;
        var min_y: i8 = 127;
        var max_x: i8 = -128;
        var max_y: i8 = -128;

        for (raw, 0..) |p, i| {
            var x = p.x;
            var y = p.y;

            if (flip) x = -x;

            for (0..rot) |_| {
                const temp = x;
                x = -y;
                y = temp;
            }

            points[i] = .{ .x = x, .y = y };

            if (x < min_x) min_x = x;
            if (y < min_y) min_y = y;
            if (x > max_x) max_x = x;
            if (y > max_y) max_y = y;
        }

        for (points) |*p| {
            p.x -= min_x;
            p.y -= min_y;
        }
        std.mem.sort(Point, points, {}, Point.less);

        return .{
            .points = points,
            .width = @as(u8, @intCast(max_x - min_x + 1)),
            .height = @as(u8, @intCast(max_y - min_y + 1)),
        };
    }

    fn isDuplicate(variants: []const Variant, target: Variant) bool {
        for (variants) |v| {
            if (v.width != target.width or v.height != target.height) continue;
            if (v.points.len != target.points.len) continue;

            var match = true;
            for (v.points, 0..) |p, i| {
                if (p.x != target.points[i].x or p.y != target.points[i].y) {
                    match = false;
                    break;
                }
            }
            if (match) return true;
        }
        return false;
    }
};

const Parser = struct {
    pub fn isQuery(line: []const u8) bool {
        return std.mem.indexOfScalar(u8, line, 'x') != null;
    }

    pub fn isHeader(line: []const u8) bool {
        return std.mem.endsWith(u8, line, ":");
    }

    pub fn parseShapeRow(alloc: std.mem.Allocator, buffer: *std.ArrayList(Point), line: []const u8, y: i8) !void {
        for (line, 0..) |c, x| {
            if (c == '#') {
                try buffer.append(alloc, .{ .x = @intCast(x), .y = y });
            }
        }
    }

    pub fn parseQuery(alloc: std.mem.Allocator, line: []const u8, shapes: []const Shape) !Query {
        var parts = std.mem.tokenizeScalar(u8, line, ':');
        const dim_part = parts.next().?;
        const count_part = parts.next().?;

        var dims = std.mem.tokenizeScalar(u8, dim_part, 'x');
        const w = try std.fmt.parseInt(u8, dims.next().?, 10);
        const h = try std.fmt.parseInt(u8, dims.next().?, 10);

        var list = std.ArrayList(usize){};
        var total_area: usize = 0;

        var tokens = std.mem.tokenizeScalar(u8, count_part, ' ');
        var shape_idx: usize = 0;
        while (tokens.next()) |t| : (shape_idx += 1) {
            const count = try std.fmt.parseInt(usize, t, 10);
            for (0..count) |_| {
                try list.append(alloc, shape_idx);
                total_area += shapes[shape_idx].area;
            }
        }

        const Context = struct {
            shapes: []const Shape,
            pub fn less(ctx: @This(), lhs: usize, rhs: usize) bool {
                return ctx.shapes[lhs].area > ctx.shapes[rhs].area;
            }
        };
        std.mem.sort(usize, list.items, Context{ .shapes = shapes }, Context.less);

        return .{
            .w = w,
            .h = h,
            .items = try list.toOwnedSlice(alloc),
            .total_area = total_area,
        };
    }
};
