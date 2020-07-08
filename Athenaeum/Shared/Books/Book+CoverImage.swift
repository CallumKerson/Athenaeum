/**
 Book+CoverImage.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

extension Book {
    var image: Image? {
        guard let coverData = getCover() else { return nil }
        let image = NSImage(data: coverData)!
        return Image(nsImage: image)
    }
}
