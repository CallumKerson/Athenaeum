/**
 AVAsset+Getters.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
import Foundation

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
                var names = artistString.components(separatedBy: " ")
                if let lastName = names.last {
                    names.removeLast(1)
                    authors
                        .append(Author(firstName: names.joined(separator: " "),
                                       lastName: lastName))
                } else {
                    authors
                        .append(Author(firstName: nil, lastName: artistString))
                }
            }
        }

        return authors
    }
}
