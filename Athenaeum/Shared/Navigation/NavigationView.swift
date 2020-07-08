/**
 NavigationView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AppNavigationView: View {
    #if os(iOS)
        @Environment(\.horizontalSizeClass) private var horizontalSizeClass
    #endif

    @ViewBuilder var body: some View {
//        #if os(iOS)
//        if horizontalSizeClass == .compact {
//            AppTabNavigation()
//        } else {
//            AppSidebarNavigation()
//        }
//        #else
        SidebarNavigationView()
            .frame(minWidth: 900, maxWidth: .infinity, minHeight: 500, maxHeight: .infinity)
    }
}

struct AppNavigationView_Previews: PreviewProvider {
    static var previews: some View {
        AppNavigationView()
    }
}
