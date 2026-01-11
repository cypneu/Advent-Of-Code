const std = @import("std");

const input = @embedFile("01.txt");

pub fn main() !void {
    var position: i32 = 50;
    var exact_stops: u32 = 0;
    var crossings: u32 = 0;

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");

    while (lines.next()) |line| {
        const is_left = line[0] == 'L';
        const amount = try std.fmt.parseInt(u32, line[1..], 10);

        const dist_to_boundary: u32 = if (is_left)
            (if (position == 0) 100 else @intCast(position))
        else
            (100 - @as(u32, @intCast(position)));

        if (amount >= dist_to_boundary) {
            crossings += 1 + (amount - dist_to_boundary) / 100;
        }

        const move = @as(i32, @intCast(amount));
        if (is_left) position -= move else position += move;

        position = @mod(position, 100);
        if (position == 0) exact_stops += 1;
    }

    std.debug.print("Part 1: {}\n", .{exact_stops});
    std.debug.print("Part 2: {}\n", .{crossings});
}
