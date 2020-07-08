/**
 Book+CoverImage.swift
 Copyright (c) 2020 Callum Kerr-Edwards */

import SwiftUI

extension Book {
    var image: Image {
        guard let coverData = getCover() else {
            return Image("cover-default")
        }
        let image = NSImage(data: coverData)!
        return Image(nsImage: image)
    }
}
