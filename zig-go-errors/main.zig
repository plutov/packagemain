const std = @import("std");

fn get_args_count(allocator: std.mem.Allocator) !usize {
    const args = try std.process.argsAlloc(allocator);
    if (args.len < 1) {
        return error.EmptyArgs;
    }

    return args.len;
}

fn get_args_count_msg(allocator: std.mem.Allocator) ![]const u8 {
    const args = try std.process.argsAlloc(allocator);
    const formatted = try std.fmt.allocPrint(allocator, "args: {d}", .{args.len});
    return formatted;
}

pub fn main() void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // const args_count = get_args_count(allocator) catch |err| {
    //     std.debug.print("invalid input: {}\n", .{err});
    //     return;
    // };

    // if (get_args_count(allocator)) |args_count| {
    //     std.debug.print("got {d} args\n", .{args_count});
    // } else |err| switch (err) {
    //     error.EmptyArgs => {
    //         std.debug.print("invalid input: no args\n", .{});
    //     },
    //     else => {
    //         std.debug.print("unexpected error: {}\n", .{err});
    //     },
    // }

    const args_count = get_args_count(allocator) catch blk: {
        break :blk 0;
    };
    std.debug.print("got {d} args\n", .{args_count});

    const msg = get_args_count_msg(allocator) catch unreachable;
    std.debug.print("{s}\n", .{msg});
}
