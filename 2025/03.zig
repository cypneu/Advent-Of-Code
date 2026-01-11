const std = @import("std");

const input = @embedFile("03.txt");

pub fn main() !void {
    var part1: u64 = 0;
    var part2: u64 = 0;

    var lines = std.mem.tokenizeAny(u8, input, "\r\n");

    while (lines.next()) |bank| {
        part1 += findMaxJoltageForTwoBatteries(bank);
        part2 += findMaxJoltageForKBatteries(bank, 12);
    }

    std.debug.print("Part 1: {d}\n", .{part1});
    std.debug.print("Part 2: {d}\n", .{part2});
}

fn findMaxJoltageForTwoBatteries(bank: []const u8) u64 {
    var max_val: u64 = 0;
    var max_seen: u8 = 0;

    for (bank, 0..) |char, i| {
        const digit = char - '0';

        if (i > 0) {
            const current_pair = @as(u64, max_seen) * 10 + digit;
            if (current_pair > max_val) max_val = current_pair;
        }

        if (digit > max_seen) max_seen = digit;
    }
    return max_val;
}

fn findMaxJoltageForKBatteries(bank: []const u8, k: usize) u64 {
    var stack: [16]u8 = undefined;
    var len: usize = 0;

    for (bank, 0..) |char, i| {
        const digit = char - '0';
        const remaining = bank.len - 1 - i;

        while (len > 0) {
            const top = stack[len - 1];
            if (digit > top and len + remaining >= k) {
                len -= 1;
            } else {
                break;
            }
        }

        if (len < k) {
            stack[len] = digit;
            len += 1;
        }
    }

    var result: u64 = 0;
    for (stack[0..len]) |d| {
        result = result * 10 + d;
    }
    return result;
}
