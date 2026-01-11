const std = @import("std");

const input = @embedFile("06.txt");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const p1 = try solvePart1(allocator);
    std.debug.print("Part 1: {d}\n", .{p1});

    const p2 = try solvePart2(allocator);
    std.debug.print("Part 2: {d}\n", .{p2});
}

fn solvePart1(allocator: std.mem.Allocator) !u64 {
    var lines = std.mem.splitBackwardsAny(u8, input, "\r\n");
    _ = lines.next();

    var operators = std.ArrayList(u8){};
    defer operators.deinit(allocator);
    const operators_line = lines.next().?;

    var it = std.mem.tokenizeScalar(u8, operators_line, ' ');
    while (it.next()) |op_str| {
        try operators.append(allocator, op_str[0]);
    }

    var result = std.ArrayList(u64){};
    try result.resize(allocator, operators.items.len);

    defer result.deinit(allocator);

    for (result.items, 0..) |*value, index| {
        switch (operators.items[index]) {
            '+' => value.* = 0,
            '*' => value.* = 1,
            else => unreachable,
        }
    }

    while (lines.next()) |line| {
        var numbers = std.mem.tokenizeScalar(u8, line, ' ');

        var index: usize = 0;
        while (numbers.next()) |number| : (index += 1) {
            const operator = operators.items[index];
            const num = try std.fmt.parseInt(u64, number, 10);
            if (operator == '+') result.items[index] += num else result.items[index] *= num;
        }
    }

    var part1: u64 = 0;
    for (result.items) |value| part1 += value;
    return part1;
}

fn solvePart2(allocator: std.mem.Allocator) !u64 {
    var grid = std.ArrayList([]const u8){};
    defer grid.deinit(allocator);

    var lines_it = std.mem.tokenizeAny(u8, input, "\r\n");
    while (lines_it.next()) |line| {
        try grid.append(allocator, line);
    }

    const height = grid.items.len;
    const width = grid.items[0].len;

    const operators_row = grid.items[height - 1];
    const number_rows = grid.items[0 .. height - 1];

    var grand_total: u64 = 0;
    var curr_sum: u64 = 0;
    var curr_prod: u64 = 1;
    var curr_op: u8 = ' ';

    var col = width;
    while (col > 0) {
        col -= 1;

        var num: u64 = 0;
        var has_digit = false;

        for (number_rows) |row| {
            const char = row[col];
            if (std.ascii.isDigit(char)) {
                has_digit = true;
                num = (num * 10) + (char - '0');
            }
        }

        if (!has_digit) {
            grand_total += if (curr_op == '*') curr_prod else curr_sum;
            curr_sum = 0;
            curr_prod = 1;
        } else {
            curr_sum += num;
            curr_prod *= num;
        }

        const op = operators_row[col];
        if (op == '+' or op == '*') curr_op = op;
    }

    grand_total += if (curr_op == '*') curr_prod else curr_sum;
    return grand_total;
}
