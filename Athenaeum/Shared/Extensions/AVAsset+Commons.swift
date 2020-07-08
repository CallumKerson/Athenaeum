/**
 AVAsset+Commons.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation

extension AVAsset {
    var commonTitle: String? {
        if let item = AVMetadataItem.metadataItems(from: metadata,
                                                   filteredByIdentifier: .commonIdentifierTitle)
            .first {
            return item.stringValue
        }
        return nil
    }

    var commonArtist: String? {
        if let item = AVMetadataItem.metadataItems(from: metadata,
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
            .metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierArtist)

        for item in items {
            if let artistString = item.stringValue {
                authors.append(Author(stringLiteral: artistString))
            }
        }

        return authors
    }
}
