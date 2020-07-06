/**
 Series.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

struct Series: Codable, Equatable, Hashable {
    public let title: String
    public let entry: Double

    public init(title: String, entry: Double) {
        self.title = title
        self.entry = entry
    }

    public var displayName: String {
        if entry.truncatingRemainder(dividingBy: 1.0) == 0.0 {
            return "Book \(Int(entry)) of \(title)"
        } else {
            return "Book \(entry) of \(title)"
        }
    }
}
