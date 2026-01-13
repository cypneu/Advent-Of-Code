const std = @import("std");

const input = @embedFile("11.txt");

const Graph = struct {
    map: std.StringHashMap([]const []const u8),
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator) Graph {
        return .{
            .map = std.StringHashMap([]const []const u8).init(allocator),
            .allocator = allocator,
        };
    }

    pub fn parse(self: *Graph, raw_data: []const u8) !void {
        var lines = std.mem.tokenizeAny(u8, raw_data, "\r\n");
        while (lines.next()) |line| {
            var parts = std.mem.splitSequence(u8, line, ": ");
            const node = parts.next() orelse continue;
            const neighbors_str = parts.next() orelse continue;

            var neighbors = std.ArrayList([]const u8){};
            var it = std.mem.splitScalar(u8, neighbors_str, ' ');
            while (it.next()) |neighbor| {
                try neighbors.append(self.allocator, neighbor);
            }

            try self.map.put(node, try neighbors.toOwnedSlice(self.allocator));
        }
    }

    pub fn countTotalPaths(self: *Graph, start: []const u8, target: []const u8) !u128 {
        var memo = std.StringHashMap(u128).init(self.allocator);
        defer memo.deinit();
        return try self.countRecursive(start, target, &memo);
    }

    fn countRecursive(self: *Graph, current: []const u8, target: []const u8, memo: *std.StringHashMap(u128)) !u128 {
        if (std.mem.eql(u8, current, target)) return 1;
        if (memo.get(current)) |count| return count;

        var total: u128 = 0;
        const neighbors = self.map.get(current) orelse return 0;

        for (neighbors) |neighbor| {
            total += try self.countRecursive(neighbor, target, memo);
        }

        try memo.put(current, total);
        return total;
    }
};

pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    var graph = Graph.init(allocator);
    try graph.parse(input);

    const part1 = try graph.countTotalPaths("you", "out");
    std.debug.print("Part 1: {d}\n", .{part1});

    const s_f = try graph.countTotalPaths("svr", "fft");
    const f_d = try graph.countTotalPaths("fft", "dac");
    const d_o = try graph.countTotalPaths("dac", "out");

    const s_d = try graph.countTotalPaths("svr", "dac");
    const d_f = try graph.countTotalPaths("dac", "fft");
    const f_o = try graph.countTotalPaths("fft", "out");

    const route_a = s_f * f_d * d_o;
    const route_b = s_d * d_f * f_o;

    std.debug.print("Part 2: {d}\n", .{route_a + route_b});
}

