const std = @import("std");

const input = @embedFile("07.txt");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    var grid = std.ArrayList([]const u8){};
    defer grid.deinit(allocator);

    var start_col: usize = 0;

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");
    while (lines.next()) |line| {
        start_col = (line.len - 1) / 2;
        try grid.append(allocator, line);
    }

    var part1: u64 = 0;
    var part2: u64 = 0;

    var current_state = std.AutoHashMap(usize, u64).init(allocator);
    defer current_state.deinit();

    var next_state = std.AutoHashMap(usize, u64).init(allocator);
    defer next_state.deinit();

    try current_state.put(start_col, 1);
    for (grid.items[1..]) |row| {
        next_state.clearRetainingCapacity();

        var it = current_state.iterator();
        while (it.next()) |entry| {
            const col = entry.key_ptr.*;
            const count = entry.value_ptr.*;

            if (row[col] == '^') {
                part1 += 1;
                if (col > 0) try addCount(&next_state, col - 1, count);
                if (col + 1 < row.len) try addCount(&next_state, col + 1, count);
            } else try addCount(&next_state, col, count);
        }

        std.mem.swap(std.AutoHashMap(usize, u64), &current_state, &next_state);
        if (current_state.count() == 0) break;
    }

    var sum_it = current_state.valueIterator();
    while (sum_it.next()) |val| part2 += val.*;

    std.debug.print("Part 1: {d}\n", .{part1});
    std.debug.print("Part 2: {d}\n", .{part2});
}

fn addCount(map: *std.AutoHashMap(usize, u64), col: usize, count: u64) !void {
    const res = try map.getOrPut(col);
    if (res.found_existing) {
        res.value_ptr.* += count;
    } else {
        res.value_ptr.* = count;
    }
}
