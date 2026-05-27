import { Component, OnInit, ChangeDetectorRef } from '@angular/core'; // 🌟 Added ChangeDetectorRef
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';

interface TicketResponse {
  ticket_number: string;
  issued_at: string;
  status: string;
}

interface CurrentQueueResponse {
  current_number: string;
}

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  currentPage: string = 'IT-05-1'; 
  queueCode: string = '--';
  timestamp: string | null = null;
  apiUrl = 'http://localhost:3000/api'; 

  // 🌟 Injected ChangeDetectorRef (cdr) into constructor to control view updates
  constructor(private http: HttpClient, private cdr: ChangeDetectorRef) {}

  ngOnInit() {
    // Initialize component lifecycle hook
  }

  // IT 05-1: Request next running ticket item
  generateQueue() {
    this.http.post<TicketResponse>(`${this.apiUrl}/tickets/next`, {}).subscribe({
      next: (res) => {
        this.queueCode = res.ticket_number;
        this.timestamp = res.issued_at;
        this.currentPage = 'IT-05-2'; 
        
        // 🌟 Force Angular to immediately re-render the HTML view
        this.cdr.detectChanges(); 
      },
      error: (err) => {
        alert('Queue request failed or conflict occurred. Please try again.');
        this.cdr.detectChanges();
      }
    });
  }

  // IT 05-1: Clear button behavior -> Redirects to preview state without resetting DB yet
  goToClearPage() {
    this.http.get<CurrentQueueResponse>(`${this.apiUrl}/queue/current`).subscribe({
      next: (res) => {
        this.queueCode = res.current_number;
        this.currentPage = 'IT-05-3'; 
        
        // 🌟 Force Angular to immediately re-render the HTML view
        this.cdr.detectChanges(); 
      },
      error: (err) => {
        this.queueCode = 'ERR';
        this.currentPage = 'IT-05-3';
        this.cdr.detectChanges();
      }
    });
  }

  // IT 05-3: Confirms and executes full data wipe cycle across backend endpoints
  executeClearQueue() {
    this.http.post<any>(`${this.apiUrl}/queue/clear`, {}).subscribe({
      next: (res) => {
        this.queueCode = '00'; 
        alert('System queue has been successfully cleared.');
        this.cdr.detectChanges(); 
      },
      error: (err) => {
        alert('Failed to execute backend queue clearing routine.');
        this.cdr.detectChanges();
      }
    });
  }

  // IT 05-2 & IT 05-3: Return route fallback
  goBackToMain() {
    this.currentPage = 'IT-05-1';
    this.cdr.detectChanges(); // 🌟 Force view update on back navigation
  }
}
