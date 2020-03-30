/**
 AVAsset+Getters.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
import Foundation
import GoodReadsKit

extension AVAsset {
    var commonTitle: String? {
        if let item = AVMetadataItem.metadataItems(from: self.metadata,
                                                   filteredByIdentifier: .commonIdentifierTitle)
            .first {
            return item.stringValue
        }
        return nil
    }

    var commonArtist: String? {
        if let item = AVMetadataItem.metadataItems(from: self.metadata,
                                                   filteredByIdentifier: .commonIdentifierArtist)
            .first {
            return item.stringValue
        }
        return nil
    }

    var commonCreationDate: String? {
        if let item = AVMetadataItem.metadataItems(from: metadata,
                                                   filteredByIdentifier: .commonIdentifierCreationDate)
            .first {
            return item.stringValue
        }
        return nil
    }

    var artistsAsAuthors: [Author] {
        var authors: [Author] = []
        let items = AVMetadataItem
            .metadataItems(from: self.metadata, filteredByIdentifier: .commonIdentifierArtist)

        for item in items {
            if let artistString = item.stringValue {
                authors.append(Author(fullName: artistString))
            }
        }

        return authors
    }
}
