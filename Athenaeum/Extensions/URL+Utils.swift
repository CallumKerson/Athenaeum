/**
 URL+Utils.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import CryptoKit
import Foundation

extension URL {
    var deSandboxedPath: String {
        path
            .replacingOccurrences(of: "Library/Containers/com.umbra.Athenaeum/Data/",
                                  with: "")
    }
}

extension URL {
    var sha256HashOfContents: String {
        do {
            return SHA256.hash(data: try Data(contentsOf: self)).description
        } catch {
            log.error("Cannot hash contents of file \(path)")
            return ""
        }
    }
}
