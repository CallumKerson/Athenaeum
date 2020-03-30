/**
 🎧📖.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
import Foundation
import GoodReadsKit

struct AudioBook: Equatable, Codable, Hashable, CustomDebugStringConvertible {
    let id: UUID
    var location: URL
    var metadata: BookMetadata?

    var debugDescription: String {
        if let title = self.metadata?.title {
            if let author = self.metadata?.authors?.author {
                return "\(title) by \(author) (\(self.location.path))"
            } else {
                return "\(title) (\(self.location.path))"
            }
        } else {
            return self.location.path
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

extension Array where Element == AudioBook {
    func hasAudibookWithSameFileAs(_ fileURL: URL) -> Bool {
        let filesInDirectory = self.map { $0.location }
        let fileSizes: [UInt64] = filesInDirectory.map { $0.fileSize }
        if fileSizes.contains(fileURL.fileSize) {
            let filesWithSameSize = filesInDirectory.filter { $0.fileSize == fileURL.fileSize }

            let fileDigest = fileURL.sha256HashOfContents
            for fileToTest in filesWithSameSize {
                if fileToTest.sha256HashOfContents == fileDigest {
                    return true
                }
            }
        }

        return false
    }
}
