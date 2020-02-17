/**
 Audiobook.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation
import RealmSwift

class Audiobook: Object, Identifiable {
    @objc dynamic let id = UUID().uuidString
    @objc dynamic var title: String
    @objc dynamic var author: String
    @objc dynamic var date: String?
    @objc dynamic var cover: Data?
    @objc dynamic var location: String?

    public required init() {
        title = "title"
        author = "author"
        super.init()
    }

    static func getBookFromFile(fileURL: URL) -> Audiobook {
        let asset = AVURLAsset(url: fileURL)

        let metadata = asset.commonMetadata

        var title: String?
        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierTitle).first {
            title = item.stringValue
        }

        var author: String?
        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierArtist).first {
            author = item.stringValue
        }

        var image: Data?
        if let artworkItem = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierArtwork).first {
            // Coerce the value to an NSData using its dataValue property
            if let imageData = artworkItem.dataValue {
                image = imageData
            }
        }

        var date: String?
        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierCreationDate).first {
            date = item.stringValue
        }

        let book = Audiobook()
        book.title = title!
        book.author = author!
        book.date = date
        book.cover = image
        book.location = fileURL.path
        return book
    }

    static func getBookFromFile(path: String) -> Audiobook {
        Audiobook.getBookFromFile(fileURL: URL(fileURLWithPath: path))
    }
}

enum AudiobookError: Error {
    case fileMissingMetadata(filepath: String)
}
