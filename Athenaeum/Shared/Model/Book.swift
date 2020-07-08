/**
 Book.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
import Foundation

struct Book: Identifiable, Equatable, Codable, Hashable {
    let id: UUID
    var audio: URL?
    var ebook: URL?
    var metadata: BookMetadata

    init(metadata: BookMetadata) {
        self.id = UUID()
        self.metadata = metadata
    }
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
