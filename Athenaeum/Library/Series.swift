/**
 Series.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

class Series {
    var title: String
    var author: String
    var entries: Set<Audiobook>

    init(title: String, author: String) {
        self.title = title
        self.author = author
        entries = Set()
    }
}

extension Series: Equatable {
    static func == (lhs: Series, rhs: Series) -> Bool {
        guard lhs.title == rhs.title else {
            return false
        }
        guard lhs.author == rhs.author else {
            return false
        }
        return true
    }
}

extension Series: Hashable {
    func hash(into hasher: inout Hasher) {
        hasher.combine(title)
        hasher.combine(author)
    }
}

extension Set where Element == Series {
    mutating func getSeries(title: String?, author: String?) -> Series? {
        if let seriesTitle = title, let seriesAuthor = author {
            if let existingSeries = first(where: { $0.title == seriesTitle && $0.author == seriesAuthor }) {
                return existingSeries
            } else {
                let newSeries = Series(title: seriesTitle, author: seriesAuthor)
                insert(newSeries)
                return newSeries
            }
        }
        return nil
    }
}
