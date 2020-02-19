/**
 URL+Desandbox.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

extension URL {
    func deSandboxedPath() -> String {
        path.replacingOccurrences(of: "Library/Containers/com.umbra.Athenaeum/Data/", with: "")
    }
}

extension String {
    var isBlank: Bool {
        trimmingCharacters(in: .whitespacesAndNewlines).isEmpty
    }

    var removeIllegalCharacters: String {
        var invalidCharacters = CharacterSet(charactersIn: ":/")
        invalidCharacters.formUnion(.newlines)
        invalidCharacters.formUnion(.illegalCharacters)
        invalidCharacters.formUnion(.controlCharacters)

        return components(separatedBy: invalidCharacters)
            .joined(separator: "")
    }
}

extension Optional where Wrapped == String {
    var isBlank: Bool {
        if let unwrapped = self {
            return unwrapped.isBlank
        } else {
            return true
        }
    }
}
