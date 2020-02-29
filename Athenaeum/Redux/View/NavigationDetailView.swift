/**
 NavigationDetailView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct NavigationDetailView: View {
    @ObservedObject var viewModel: NavigationDetailViewModel

    init(_ viewModel: NavigationDetailViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        ScrollView {
            Unwrap(viewModel.audiobook) { book in
                VStack(alignment: HorizontalAlignment.leading, spacing: 12) {
                    HStack(alignment: VerticalAlignment.center, spacing: 24) {
                        Unwrap(book.getCover()) { coverData in
                            DetailCover(data: coverData)
                        }
                        BookText(audiobook: book)
                            .font(.subheadline)
                    }

                    Unwrap(book.bookDescription) { description in
                        Divider()
                        Text(description)
                            .lineLimit(nil)
                    }
                }
            }
        }
        .padding()
        .frame(minWidth: 800, maxWidth: 800, minHeight: 500)
    }
}

struct BookText: View {
    let audiobook: AudioBook

    var body: some View {
        VStack(alignment: HorizontalAlignment.leading) {
            VStack(alignment: HorizontalAlignment.leading) {
                Text(audiobook.title).font(.title)

                Unwrap(audiobook.series) { series in
                    Text("Book \(series.entry) of \(series.title)")
                }
            }.padding(.bottom)

            Unwrap(audiobook.getAuthorsString()) { author in
                Text("Written by ") + Text(author)
            }
            Unwrap(audiobook.publicationDate) { releaseDate in
                Text("Released on ") + Text(releaseDate)
            }
            Unwrap(audiobook.narrator) { narrator in
                Text("Narrated by ") + Text(narrator)
            }
            Unwrap(audiobook.isbn) { isbn in
                Text("ISBN ") + Text(isbn)
            }
        }
        .font(.subheadline)
    }
}

struct DetailCover: View {
    let data: Data

    var body: some View {
        Image(nsImage: NSImage(data: data)!)
            .resizable()
            .scaledToFit()
            .frame(width: 400, height: 400)
            .clipShape(RoundedRectangle(cornerRadius: 10))
            .shadow(radius: 10)
    }
}

#if DEBUG
    struct NavigationDetailView_Previews: PreviewProvider {
        static var previews: some View {
            NavigationDetailView(NavigationDetailViewModel(id: sampleAudiobook.id,
                                                           store: sampleStore))
        }
    }
#endif