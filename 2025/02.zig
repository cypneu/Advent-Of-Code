const std = @import("std");

const input = @embedFile("02.txt");

pub fn main() !void {
    var range_tokens = std.mem.tokenizeAny(u8, input, ",\n\r ");

    var part1: u64 = 0;
    var part2: u64 = 0;

    while (range_tokens.next()) |range_str| {
        var parts = std.mem.tokenizeScalar(u8, range_str, '-');

        const start = try std.fmt.parseInt(u64, parts.next().?, 10);
        const end = try std.fmt.parseInt(u64, parts.next().?, 10) + 1;

        for (start..end) |current| {
            if (isInvalid1(current)) part1 += current;
            if (isInvalid2(current)) part2 += current;
        }
    }

    std.debug.print("Part 1: {}\n", .{part1});
    std.debug.print("Part 2: {}\n", .{part2});
}

fn isInvalid1(n: u64) bool {
    const digits = std.math.log10_int(n) + 1;
    if (digits % 2 != 0) return false;

    const half = digits / 2;
    const divisor = std.math.pow(u64, 10, half);

    return (n / divisor) == (n % divisor);
}

fn isInvalid2(n: u64) bool {
    const len = std.math.log10_int(n) + 1;

    for (1..len / 2 + 1) |p| {
        if (len % p != 0) continue;

        const mask = std.math.pow(u64, 10, p);
        var target = n / mask;

        while (target > 0) : (target /= mask) {
            if (target % mask != n % mask) break;
        } else return true;
    }

    return false;
}
