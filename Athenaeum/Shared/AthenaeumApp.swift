/**
 AthenaeumApp.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

@main
struct AthenaeumApp: App {
    @StateObject private var model = AthenaeumModel()

    var body: some Scene {
        WindowGroup {
            AppNavigationView().environmentObject(model)
        }
    }
}

struct AthenaeumApp_Previews: PreviewProvider {
    static var previews: some View {
        /*@START_MENU_TOKEN@*/Text("Hello, World!")/*@END_MENU_TOKEN@*/
    }
}
