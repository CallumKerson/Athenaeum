/**
 AudiobookCellView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct AudiobookCellView: View {
    let book: Audiobook
    var body: some View {
        HStack {
            Unwrap(book.getCover()) { coverData in
                SmallCover(data: coverData)
            }
            VStack(alignment: .leading) {
                Text(book.title).font(.subheadline)

                Text(book.author).font(.footnote)
            }
        }
        .padding(5)
    }
}

struct SmallCover: View {
    let data: Data
    let color = Color(red: 232 / 255, green: 238 / 255, blue: 246 / 255)

    var body: some View {
        let image = NSImage(data: data)!
        return Image(nsImage: image.resize(withSize: NSMakeSize(50, 50))!)
            .frame(width: 50, height: 50)
            .clipShape(RoundedRectangle(cornerRadius: 5))
            .shadow(radius: 3)
    }
}

#if DEBUG
    struct AudiobookCellView_Previews: PreviewProvider {
        static var previews: some View {
            Group {
                AudiobookCellView(book: previewAudiobooks[0])
                AudiobookCellView(book: previewAudiobooks[1])
                AudiobookCellView(book: previewAudiobooks[2])
            }
        }
    }
#endif
