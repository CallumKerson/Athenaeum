/**
 Audiobook+Construction.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation
import GoodReadsKit

extension AudiobookFile {
    convenience init(fromFileWithPath path: String, withGoodReads goodReads: GoodReads? = PreferencesStore.global.goodReadsAPI) {
        self.init(fromFile: URL(fileURLWithPath: path), withGoodReads: goodReads)
    }

    convenience init(fromFile fileURL: URL, withGoodReads goodReads: GoodReads? = PreferencesStore.global.goodReadsAPI) {
        log.debug("Creating Audiobook from file \(fileURL.path)")
        let asset = AVURLAsset(url: fileURL)

        let metadata = asset.commonMetadata

        var title: String?
        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierTitle).first {
            title = item.stringValue!.removeIllegalCharacters
        }

        var author: String?
        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierArtist).first {
            author = item.stringValue!.removeIllegalCharacters
        }

        if let goodReads = goodReads {
            log.debug("Getting audiobook metadata from GoodReads API")

            do {
                let fetchedBook = try goodReads
                    .getBook(title: title!, author: author!)
                var series: Series?
                if let seriesTitle = fetchedBook.seriesTitle, let seriesEntry = fetchedBook.seriesEntry {
                    series = Series(title: seriesTitle, entry: String(seriesEntry))
                }
                self.init(title: fetchedBook.title,
                          author: fetchedBook.getAuthorString(),
                          file: fileURL,
                          publicationDate: fetchedBook.getDateString(),
                          isbn: fetchedBook.isbn,
                          summary: fetchedBook.bookDescription,
                          series: series)
                return
            } catch {
                log.error("Could not get book details from GoodReads API with search terms title: \(title ?? "nil") and author \(author ?? "nil")")
                log.error(error)
            }
        }

        var date: String?
        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierCreationDate).first {
            date = item.stringValue
        }
        log.debug("Getting Audiobook metadata from file metadata")
        self.init(title: title!, author: author!, file: fileURL, publicationDate: date)
    }
}

extension AudiobookFile {
    func getCover() -> Data? {
        if let artworkItem = AVMetadataItem.metadataItems(from: AVURLAsset(url: location).commonMetadata, filteredByIdentifier: .commonIdentifierArtwork).first {
            // Coerce the value to an NSData using its dataValue property
            if let imageData = artworkItem.dataValue {
                return imageData
            }
        }
        return nil
    }
}

enum AudiobookFileError: Error {
    case fileMissingMetadata(filepath: String)
}

extension PreferencesStore {
    var goodReadsAPI: GoodReads? {
        if goodReadsAPIKey.isBlank {
            return nil
        }
        return GoodReadsRESTAPI(apiKey: goodReadsAPIKey)
    }
}
