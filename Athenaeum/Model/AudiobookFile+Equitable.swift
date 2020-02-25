/**
 AudiobookFile+Equitable.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import CryptoKit
import Foundation

extension AudiobookFile: Equatable {
    static func == (lhs: AudiobookFile, rhs: AudiobookFile) -> Bool {
        guard lhs.title == rhs.title else {
            return false
        }
        guard lhs.author == rhs.author else {
            return false
        }
        guard lhs.location == rhs.location else {
            return false
        }
        return true
    }
}

extension AudiobookFile: Hashable {
    func hash(into hasher: inout Hasher) {
        hasher.combine(title)
        hasher.combine(author)
        hasher.combine(location)
    }
}
