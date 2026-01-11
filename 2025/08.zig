const std = @import("std");

const input = @embedFile("08.txt");

const JunctionBox = struct {
    x: u64,
    y: u64,
    z: u64,

    pub fn distanceSq(self: JunctionBox, other: JunctionBox) u64 {
        const dx = if (self.x > other.x) self.x - other.x else other.x - self.x;
        const dy = if (self.y > other.y) self.y - other.y else other.y - self.y;
        const dz = if (self.z > other.z) self.z - other.z else other.z - self.z;
        return dx * dx + dy * dy + dz * dz;
    }
};

const Edge = struct {
    u: usize,
    v: usize,
    dist_sq: u64,

    pub fn lessThan(_: void, a: Edge, b: Edge) bool {
        return a.dist_sq < b.dist_sq;
    }
};

const UnionFind = struct {
    parent: []usize,
    size: []u64,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, n: usize) !UnionFind {
        const parent = try allocator.alloc(usize, n);
        const size = try allocator.alloc(u64, n);

        for (0..n) |i| {
            parent[i] = i;
            size[i] = 1;
        }

        return UnionFind{
            .parent = parent,
            .size = size,
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *UnionFind) void {
        self.allocator.free(self.parent);
        self.allocator.free(self.size);
    }

    pub fn find(self: *UnionFind, i: usize) usize {
        var root = i;
        while (root != self.parent[root]) {
            self.parent[root] = self.parent[self.parent[root]];
            root = self.parent[root];
        }
        return root;
    }

    pub fn unionSets(self: *UnionFind, i: usize, j: usize) bool {
        const root_i = self.find(i);
        const root_j = self.find(j);

        if (root_i != root_j) {
            if (self.size[root_i] < self.size[root_j]) {
                self.parent[root_i] = root_j;
                self.size[root_j] += self.size[root_i];
            } else {
                self.parent[root_j] = root_i;
                self.size[root_i] += self.size[root_j];
            }
            return true;
        }
        return false;
    }
};

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    var junction_boxes = std.ArrayList(JunctionBox){};
    defer junction_boxes.deinit(allocator);

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");
    while (lines.next()) |line| {
        var coords = std.mem.splitScalar(u8, line, ',');
        try junction_boxes.append(allocator, .{
            .x = try std.fmt.parseInt(u64, coords.next().?, 10),
            .y = try std.fmt.parseInt(u64, coords.next().?, 10),
            .z = try std.fmt.parseInt(u64, coords.next().?, 10),
        });
    }

    var edges = std.ArrayList(Edge){};
    defer edges.deinit(allocator);

    const items = junction_boxes.items;
    for (0..items.len) |i| {
        for (i + 1..items.len) |j| {
            try edges.append(allocator, .{
                .u = i,
                .v = j,
                .dist_sq = items[i].distanceSq(items[j]),
            });
        }
    }

    std.mem.sort(Edge, edges.items, {}, Edge.lessThan);

    var uf = try UnionFind.init(allocator, items.len);
    defer uf.deinit();

    var num_components = items.len;
    const p1_limit = @min(1000, edges.items.len);

    for (edges.items[0..p1_limit]) |edge| {
        if (uf.unionSets(edge.u, edge.v)) {
            num_components -= 1;
        }
    }

    var sizes = std.ArrayList(u64){};
    for (0..uf.parent.len) |i| {
        if (uf.parent[i] == i) {
            try sizes.append(allocator, uf.size[i]);
        }
    }
    defer sizes.deinit(allocator);

    std.mem.sort(u64, sizes.items, {}, std.sort.desc(u64));

    var part1: u64 = 1;
    const count = @min(3, sizes.items.len);
    for (0..count) |k| part1 *= sizes.items[k];

    std.debug.print("Part 1: {}\n", .{part1});

    for (edges.items[p1_limit..]) |edge| {
        if (uf.unionSets(edge.u, edge.v)) {
            num_components -= 1;

            if (num_components == 1) {
                const b1 = items[edge.u];
                const b2 = items[edge.v];
                std.debug.print("Part 2: {}\n", .{b1.x * b2.x});
                return;
            }
        }
    }
}
