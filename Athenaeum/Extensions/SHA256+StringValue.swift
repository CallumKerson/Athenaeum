/**
 SHA256+StringValue.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import CryptoKit
import Foundation

extension SHA256.Digest {
    var stringValue: String {
        map { String(format: "%02hhx", $0) }.joined()
    }
}
