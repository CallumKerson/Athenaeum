/**
 NSImage+Resize.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */
import Cocoa

extension NSImage {
    func resize(withSize targetSize: NSSize) -> NSImage? {
        let frame = NSRect(x: 0, y: 0, width: targetSize.width,
                           height: targetSize.height)
        guard let representation = self.bestRepresentation(for: frame,
                                                           context: nil,
                                                           hints: nil) else {
            return nil
        }
        let image = NSImage(size: targetSize, flipped: false,
                            drawingHandler: { (_) -> Bool in
                                representation.draw(in: frame)
        })

        return image
    }
}
