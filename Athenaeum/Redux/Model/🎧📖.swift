/**
 🎧📖.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import AVFoundation

struct AudioBook: Equatable, Codable, Hashable, CustomDebugStringConvertible {
    let id: UUID
    var location: URL
    var contentsHash: String
    var title: String
    var authors: [Author]?
    var narrator: String?
    var publicationDate: String?
    var isbn: String?
    var bookDescription: String?
    var series: Series?

    public func getAuthorsString() -> String? {
        if let authors = authors {
            let authorsStrings = authors.map { $0.getAuthorString() }
            return authorsStrings.joined(separator: ", ")
                .replacingLastOccurrenceOfString(",", with: " &")
        } else {
            return nil
        }
    }
    
    var debugDescription: String {
        if let authors = getAuthorsString() {
            return "\(title) by \(authors) (\(location.path))"
        } else {
            return "\(title) (\(location.path))"
        }
    }
}

struct Author: Equatable, Codable, Hashable {
    let firstName: String?
    let lastName: String

    func getAuthorString() -> String {
        if let firstName = firstName {
            return "\(firstName) \(self.lastName)"
        } else {
            return self.lastName
        }
    }
}


extension AudioBook {
    func getCover() -> Data? {
        if let artworkItem = AVMetadataItem
            .metadataItems(from: AVURLAsset(url: location).commonMetadata,
                           filteredByIdentifier: .commonIdentifierArtwork)
            .first {
            // Coerce the value to an NSData using its dataValue property
            if let imageData = artworkItem.dataValue {
                return imageData
            }
        }
        return nil
    }
}
