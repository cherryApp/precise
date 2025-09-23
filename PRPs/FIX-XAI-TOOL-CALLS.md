Some x.ai tool calls runs into a failer.

<xai:function_call name="multiedit"> [{"new_string":"import { Component, OnInit, signal, computed, inject } from
'@angular/core';\nimport { CommonModule } from '@angular/common';\nimport { FormsModule } from
'@angular/forms';\nimport { firstValueFrom } from 'rxjs';\nimport { SiteConfigService } from '../../services';\nimport
{ ApiService } from '../../services';\nimport { NotificationService } from '../../services';, "

• find the application part that is responsible for tool calls
• update the functionality to a more flexible tool call handling and prevent this errors