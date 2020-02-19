/**
 AudiobookCellView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AudiobookCellView: View {
    let book: Audiobook
    var body: some View {
        HStack {
            if book.getCover() != nil {
                SmallCover(data: book.getCover()!)
            }
            VStack(alignment: .leading) {
                Text(book.title).font(.subheadline)

                Text(book.author).font(.footnote)
            }
        }
    }
}

struct SmallCover: View {
    let data: Data

    var body: some View {
        Image(nsImage: NSImage(data: data)!)
            .resizable()
            .scaledToFit()
            .aspectRatio(contentMode: .fit)
            .frame(width: 50, height: 50).clipShape(RoundedRectangle(cornerRadius: 5))
            .shadow(radius: 10)
    }
}

struct AudiobookCellView_Previews: PreviewProvider {
    static var previews: some View {
        AudiobookCellView(book: Audiobook(fromFileWithPath: "/Users/ckerson/Music/TWoK.m4b"))
    }
}
