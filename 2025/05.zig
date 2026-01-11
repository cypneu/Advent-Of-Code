const std = @import("std");

const input = @embedFile("05.txt");

const EventType = enum(u8) {
    range_start,
    point,
    range_end,
};

const Event = struct {
    num: u64,
    type: EventType,

    pub fn lessThan(_: void, lhs: Event, rhs: Event) bool {
        if (lhs.num != rhs.num) return lhs.num < rhs.num;
        return @intFromEnum(lhs.type) < @intFromEnum(rhs.type);
    }
};

const Interval = struct {
    start: u64,
    end: u64,

    pub fn lessThan(_: void, lhs: Interval, rhs: Interval) bool {
        if (lhs.start != rhs.start) return lhs.start < rhs.start;
        return lhs.end < rhs.end;
    }
};

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    var events = std.ArrayList(Event){};
    defer events.deinit(allocator);

    var intervals = std.ArrayList(Interval){};
    defer intervals.deinit(allocator);

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");
    while (lines.next()) |line| {
        var it = std.mem.splitScalar(u8, line, '-');

        const first_part = it.next().?;
        const val1 = try std.fmt.parseInt(u64, first_part, 10);

        if (it.next()) |second_part| {
            const val2 = try std.fmt.parseInt(u64, second_part, 10);
            try events.append(allocator, .{ .num = val1, .type = .range_start });
            try events.append(allocator, .{ .num = val2, .type = .range_end });

            try intervals.append(allocator, .{ .start = val1, .end = val2 });
        } else {
            try events.append(allocator, .{ .num = val1, .type = .point });
        }
    }

    var part1: u32 = 0;
    std.sort.block(Event, events.items, {}, Event.lessThan);

    var active_ranges: u32 = 0;
    for (events.items) |event| {
        switch (event.type) {
            .range_start => active_ranges += 1,
            .point => if (active_ranges > 0) {
                part1 += 1;
            },
            .range_end => active_ranges -= 1,
        }
    }

    var part2: u64 = 0;
    std.sort.block(Interval, intervals.items, {}, Interval.lessThan);

    var current = intervals.items[0];
    for (intervals.items[1..]) |interval| {
        if (interval.start <= current.end) {
            current.end = @max(current.end, interval.end);
        } else {
            part2 += current.end - current.start + 1;
            current = interval;
        }
    }
    part2 += current.end - current.start + 1;

    std.debug.print("Part 1: {d}\n", .{part1});
    std.debug.print("Part 2: {d}\n", .{part2});
}
