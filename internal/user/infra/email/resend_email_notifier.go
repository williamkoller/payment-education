package infra_email

import (
	"fmt"

	port_email_notifier "github.com/williamkoller/system-education/internal/user/port/email"
	"github.com/williamkoller/system-education/shared/infra/email"
)

type ResendEmailNotifier struct {
	client *email.ResendClient
}

func NewResendEmailNotifier(client *email.ResendClient) port_email_notifier.EmailNotifier {
	return &ResendEmailNotifier{client: client}
}

func (n *ResendEmailNotifier) SendWelcomeEmail(name, emailAddr string) error {
	subject := "Bem-vindo á Plataforma"

	html := fmt.Sprintf(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html dir="ltr" lang="en">
  <head>
    <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
    <meta name="x-apple-disable-message-reformatting" />
    <style>
      @font-face {
        font-family: 'Inter';
        font-style: normal;
        font-weight: 400;
        mso-font-alt: 'Helvetica';
        src: url(https://fonts.gstatic.com/s/inter/v18/UcCO3FwrK3iLTeHuS_nVMrMxCp50SjIw2boKoduKmMEVuLyfAZ9hiA.woff2) format('woff2');
      }

      * {
        font-family: 'Inter', Helvetica;
      }
    </style>
    <style>
      @font-face {
        font-family: 'Inter';
        font-style: normal;
        font-weight: 600;
        mso-font-alt: 'Helvetica';
        src: url(https://fonts.gstatic.com/s/inter/v18/UcC73FwrK3iLTeHuS_fjbvMwCp50PDca1ZL7.woff2) format('woff2');
      }

      * {
        font-family: 'Inter', Helvetica;
      }
    </style>
    <style>
      @font-face {
        font-family: 'Inter';
        font-style: normal;
        font-weight: 700;
        mso-font-alt: 'Helvetica';
        src: url(https://fonts.gstatic.com/s/inter/v18/UcC73FwrK3iLTeHuS_fjbvMwCp50BTca1ZL7.woff2) format('woff2');
      }

      * {
        font-family: 'Inter', Helvetica;
      }
    </style>
  </head>
  <body style="margin:0;margin-right:0;margin-left:0">
    <!--$--><!--html--><!--head--><!--body-->
    <table
      border="0"
      width="100%%"
      cellpadding="0"
      cellspacing="0"
      role="presentation"
      align="center">
      <tbody>
        <tr>
          <td style="margin:0;margin-left:12px;margin-right:12px">
            <table
              align="center"
              width="100%%"
              border="0"
              cellpadding="0"
              cellspacing="0"
              role="presentation"
              style="max-width:37.5em;margin-left:auto;margin-right:auto;box-sizing:border-box;padding-top:1rem;padding-bottom:1rem;">
              <tbody>
                <tr style="width:100%%">
                  <td>
                    <h1>Olá, %s!</h1>
                    <p>Seja muito bem‑vindo à nossa plataforma System Education.</p>
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
      </tbody>
    </table>
    <!--/$-->
  </body>
</html>`, name)

	if err := n.client.SendEmail(emailAddr, subject, html); err != nil {
		return fmt.Errorf("failed to send welcome email to %s: %w", emailAddr, err)
	}

	return nil
}
