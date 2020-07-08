/**
 Series.swift
 Copyright (c) 2020 Callum Kerr-Edwards
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
        if self.entry.truncatingRemainder(dividingBy: 1.0) == 0.0 {
            return "Book \(Int(self.entry)) of \(self.title)"
        } else {
            return "Book \(self.entry) of \(self.title)"
        }
    }
}
