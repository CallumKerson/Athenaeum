/**
 Audiobook+GetFromFile.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import AVFoundation
import Foundation

extension Audiobook {
    func getCover() -> Data? {
        if let artworkItem = AVMetadataItem.metadataItems(from: AVURLAsset(url: file).commonMetadata, filteredByIdentifier: .commonIdentifierArtwork).first {
            // Coerce the value to an NSData using its dataValue property
            if let imageData = artworkItem.dataValue {
                return imageData
            }
        }
        return nil
    }
}
