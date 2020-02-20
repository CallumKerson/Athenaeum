/**
 Audiobook+Construction.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation
import GoodReadsKit

extension Audiobook {
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
                var series: (title: String, entry: String)?
                if let seriesTitle = fetchedBook.seriesTitle, let seriesEntry = fetchedBook.seriesEntry {
                    series = (title: seriesTitle, entry: String(seriesEntry))
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

enum AudiobookError: Error {
    case fileMissingMetadata(filepath: String)
}

extension String {
    func removeCharacters(from forbiddenChars: CharacterSet) -> String {
        let passed = unicodeScalars.filter { !forbiddenChars.contains($0) }
        return String(String.UnicodeScalarView(passed))
    }

    func removeCharacters(from: String) -> String {
        removeCharacters(from: CharacterSet(charactersIn: from))
    }
}

extension PreferencesStore {
    var goodReadsAPI: GoodReads? {
        if goodReadsAPIKey.isBlank {
            return nil
        }
        return GoodReadsRESTAPI(apiKey: goodReadsAPIKey)
    }
}
