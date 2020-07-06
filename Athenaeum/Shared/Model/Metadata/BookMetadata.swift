/**
 BookMetadata.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

struct BookMetadata: Equatable, Codable, Hashable {
    var title: String
    var authors: [Author]?
    var goodReadsID: String?
    var narrators: [Author]?
    var illustrators: [Author]?
    var series: Series?
    var isbn: String?
    var isbn13: String?
    var publicationDate: PublicationDate?
    var summary: String?
    var genre: Genre?

    init(title: String) {
        self.title = title
    }
}
