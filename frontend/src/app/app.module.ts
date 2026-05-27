import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http'; // Import the new functions
import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    // Leave empty since AppComponent is defined as standalone: true
  ],
  imports: [
    BrowserModule,
    AppComponent      // Imports our main queue component layout
  ],
  providers: [
    provideHttpClient(withInterceptorsFromDi()) // Modern way to enable HTTP Client in Angular
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
