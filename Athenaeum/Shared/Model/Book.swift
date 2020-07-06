/**
 Book.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation

struct Book: Identifiable, Equatable, Codable, Hashable {
    let id: UUID
    let audio: URL?
    let ebook: URL?
    let metadata: BookMetadata?
}

extension Book {
    func getCover() -> Data? {
        if let audio = audio {
            guard let artworkItem = AVMetadataItem
                .metadataItems(from: AVURLAsset(url: audio).commonMetadata,
                               filteredByIdentifier: .commonIdentifierArtwork)
                .first else { return nil }
            return artworkItem.dataValue
        }
        // TODO: add get cover from ebook
        return nil
    }
}
