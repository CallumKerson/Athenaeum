/**
 String+Utils.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

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

    var trimmed: String {
        trimmingCharacters(in: CharacterSet.whitespacesAndNewlines)
    }

    func removeCharacters(from forbiddenChars: CharacterSet) -> String {
        let passed = unicodeScalars.filter { !forbiddenChars.contains($0) }
        return String(String.UnicodeScalarView(passed))
    }

    func removeCharacters(from: String) -> String {
        self.removeCharacters(from: CharacterSet(charactersIn: from))
    }

    func replacingLastOccurrenceOfString(_ searchString: String,
                                         with replacementString: String) -> String {
        if let range = self.range(of: searchString,
                                  options: [.backwards],
                                  range: nil,
                                  locale: nil) {
            return replacingCharacters(in: range, with: replacementString)
        }
        return self
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
