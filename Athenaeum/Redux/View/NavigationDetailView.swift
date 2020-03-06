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
                                .padding(8)
                        }
                        VStack(alignment: HorizontalAlignment.leading) {
                            if self.viewModel.isGoodReadsConfigured {
                                HStack {
                                    Spacer()
                                    Button(action: {
                                        self.viewModel.showFixMatchDialog()
                                    }) {
                                        Text("Fix Match")
                                            .frame(maxWidth: 100, maxHeight: 24)
                                    }
                                }
                            }
                            Spacer()
                            BookText(audiobook: book)
                            Spacer()
                        }
                    }

                    Unwrap(book.bookDescription) { description in
                        Divider()
                        SummaryView(summary: description)
                            .lineLimit(nil)
                    }
                }
            }
        }
        .blur(radius: viewModel.isImporting ? 10 : 0)
        .padding()
        .overlay(
            VStack {
                if viewModel.isImporting {
                    ActivityIndicator()
                        .frame(width: 100, height: 100)
                } else {
                    EmptyView()
                }
            }
        )
        .frame(minWidth: 800, maxWidth: 800, minHeight: 500)
        .sheet(isPresented: $viewModel.fixMatchDialogDisplayed) {
            Unwrap(self.viewModel.audiobook) { audiobook in
                FixMatchModalView(viewModel: FixMatchModalViewModel(audiobook: audiobook,
                                                                    store: self.viewModel.store))
            }
        }
    }
}

struct BookText: View {
    let audiobook: AudioBook

    var body: some View {
        VStack(alignment: HorizontalAlignment.leading) {
            VStack(alignment: HorizontalAlignment.leading) {
                Unwrap(audiobook.title) { title in
                    Text(title).font(.title)
                }

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
            .shadow(radius: 5)
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
