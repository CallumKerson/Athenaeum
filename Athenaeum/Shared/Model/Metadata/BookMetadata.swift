/**
 BookMetadata.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
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

extension BookMetadata {
    static func fromAudiobook(audiobook audiobookURL: URL) -> BookMetadata? {
        let audiobookAsset = AVURLAsset(url: audiobookURL)
        guard let title = audiobookAsset.commonTitle else { return nil }

        var metadata = BookMetadata(title: title)
        metadata.authors = audiobookAsset.artistsAsAuthors

        if let date = audiobookAsset.commonCreationDate {
            metadata.publicationDate = try? PublicationDate(from: date)
        }
        return metadata
    }
}
