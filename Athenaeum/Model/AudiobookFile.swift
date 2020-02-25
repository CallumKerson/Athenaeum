/**
 AudiobookFile.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

class AudiobookFile: Audiobook {
    var title: String
    var author: String
    var location: URL
    var narrator: String?
    var publicationDate: String?
    var isbn: String?
    var summary: String?
    var series: Series?

    init(title: String,
         author: String,
         file: URL,
         narrator: String? = nil,
         publicationDate: String? = nil,
         isbn: String? = nil,
         summary: String? = nil,
         series: Series? = nil) {
        self.title = title
        self.author = author
        self.location = file
        self.narrator = narrator
        self.publicationDate = publicationDate
        self.isbn = isbn
        self.summary = summary
        self.series = series
    }
}

extension AudiobookFile: Identifiable {
//    let id = UUID()
}

extension AudiobookFile: CustomStringConvertible {
    var description: String {
        "\(self.title) by \(self.author)"
    }
}
