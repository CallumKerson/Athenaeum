/**
 AudiobookView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AudiobookView: View {
    let book: Audiobook

    var body: some View {
        VStack {
            if book.cover != nil {
                Cover(data: book.cover!)
            }
            VStack(alignment: .leading) {
                Text(book.title)
                    .font(.headline)
                Text(book.author)
                    .font(.subheadline)

                if book.date != nil {
                    Text(book.date!)
                        .font(.subheadline)
                }
            }
        }
        .padding()
    }
}

struct Cover: View {
    let data: Data

    var body: some View {
        Image(nsImage: NSImage(data: data)!)
            .resizable()
            .scaledToFit()
            .frame(width: 400, height: 400)
            .clipShape(RoundedRectangle(cornerRadius: 5))
            .shadow(radius: 10)
    }
}

#if DEBUG
    struct AudiobookView_Previews: PreviewProvider {
        static var previews: some View {
            AudiobookView(book: Audiobook.getBookFromFile(path: "/Users/ckerson/Music/TWoK.m4b"))
        }
    }
#endif
