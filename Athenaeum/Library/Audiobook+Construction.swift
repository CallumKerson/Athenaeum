/**
 Audiobook+Construction.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation
import GoodReadsKit

extension Audiobook {
    convenience init(fromFileWithPath path: String) {
        self.init(fromFile: URL(fileURLWithPath: path))
    }

    convenience init(fromFile fileURL: URL) {
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
        if !Preferences.getString(for: .goodReadsAPIKey).isBlank {
            log.debug("Getting audiobook metadata from GoodReads API")

            do {
                let fetchedBook = try GoodReads(apiKey: Preferences.getString(for: .goodReadsAPIKey)!)
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

//    static func getBookFromFile(fileURL: URL) -> Audiobook {
//        log.debug("Creating Audiobook from file \(fileURL.path)")
//        let asset = AVURLAsset(url: fileURL)
//
//        let metadata = asset.commonMetadata
//
//        var title: String?
//        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierTitle).first {
//            title = item.stringValue!.removeIllegalCharacters
//        }
//
//        var author: String?
//        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierArtist).first {
//            author = item.stringValue!.removeIllegalCharacters
//        }
//        if !Preferences.getString(for: .goodReadsAPIKey).isBlank {
//            log.debug("Getting audiobook metadata from GoodReads API")
//
//            do {
//                let fetchedBook = try GoodReads(apiKey: Preferences.getString(for: .goodReadsAPIKey)!)
//                    .getBook(title: title!, author: author!)
//                var series: (title: String, entry: String)?
//                if let seriesTitle = fetchedBook.seriesTitle, let seriesEntry = fetchedBook.seriesEntry {
//                    series = (title: seriesTitle, entry: String(seriesEntry))
//                }
//                return Audiobook(title: fetchedBook.title,
//                                 author: fetchedBook.getAuthorString(),
//                                 file: fileURL,
//                                 publicationDate: fetchedBook.getDateString(),
//                                 isbn: fetchedBook.isbn,
//                                 summary: fetchedBook.bookDescription,
//                                 series: series)
//            } catch {
//                log.error("Could not get book details from GoodReads API with search terms title: \(title ?? "nil") and author \(author ?? "nil")")
//                log.error(error)
//            }
//        }
//        var date: String?
//        if let item = AVMetadataItem.metadataItems(from: metadata, filteredByIdentifier: .commonIdentifierCreationDate).first {
//            date = item.stringValue
//        }
//        log.debug("Getting Audiobook metadata from file metadata")
//        return Audiobook(title: title!, author: author!, file: fileURL, publicationDate: date)
//    }
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
