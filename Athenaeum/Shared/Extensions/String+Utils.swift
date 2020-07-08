/**
 String+Utils.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

extension String {
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
