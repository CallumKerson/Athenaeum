/**
 Audiobook+Equitable.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import CryptoKit
import Foundation

extension Audiobook: Equatable {
    static func == (lhs: Audiobook, rhs: Audiobook) -> Bool {
        guard lhs.title == rhs.title else {
            return false
        }
        guard lhs.author == rhs.author else {
            return false
        }
        return true
    }
}

extension Audiobook: Hashable {
    func hash(into hasher: inout Hasher) {
        hasher.combine(title)
        hasher.combine(author)
    }
}
