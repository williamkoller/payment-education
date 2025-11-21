package infra_email

import (
	"fmt"
	"time"

	port_email_notifier "github.com/williamkoller/system-education/internal/user/port/email"
	"github.com/williamkoller/system-education/shared/infra/email"
)

type ResendEmailNotifier struct {
	client email.EmailClient
}

func NewResendEmailNotifier(client email.EmailClient) port_email_notifier.EmailNotifier {
	return &ResendEmailNotifier{client: client}
}

func (n *ResendEmailNotifier) SendWelcomeEmail(name, emailAddr string) error {
	subject := "Bem-vindo á Plataforma"

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-BR">
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="x-apple-disable-message-reformatting" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />

    <!-- Fontes -->
    <style>
      @font-face {
        font-family: 'Inter';
        font-style: normal;
        font-weight: 400;
        mso-font-alt: 'Helvetica';
        src: url(https://fonts.gstatic.com/s/inter/v18/UcCO3FwrK3iLTeHuS_nVMrMxCp50SjIw2boKoduKmMEVuLyfAZ9hiA.woff2) format('woff2');
      }
      @font-face {
        font-family: 'Inter';
        font-style: normal;
        font-weight: 600;
        mso-font-alt: 'Helvetica';
        src: url(https://fonts.gstatic.com/s/inter/v18/UcC73FwrK3iLTeHuS_fjbvMwCp50PDca1ZL7.woff2) format('woff2');
      }
      @font-face {
        font-family: 'Inter';
        font-style: normal;
        font-weight: 700;
        mso-font-alt: 'Helvetica';
        src: url(https://fonts.gstatic.com/s/inter/v18/UcC73FwrK3iLTeHuS_fjbvMwCp50BTca1ZL7.woff2) format('woff2');
      }
      * {
        font-family: 'Inter', Helvetica, Arial, sans-serif;
      }
    </style>
  </head>

  <body style="margin:0; padding:0; background-color:#f5f5f7;">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" align="center">
      <tr>
        <td style="padding:24px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" align="center" style="max-width:600px; background:#ffffff; border-radius:8px; padding:32px;">
            <tr>
              <td style="text-align:left;">
                
                <h1 style="font-size:24px; font-weight:700; color:#111; margin:0 0 16px 0;">
                  Olá, %s!
                </h1>

                <p style="font-size:16px; color:#444; margin:0 0 12px 0; line-height:1.5;">
                  Seja muito bem-vindo à <strong>System Education</strong>! É um prazer ter você conosco.
                </p>

                <p style="font-size:16px; color:#444; margin:0 0 12px 0; line-height:1.5;">
                  A partir de agora, você terá acesso a recursos, conteúdos e ferramentas pensadas para facilitar sua jornada de aprendizado.
                </p>

                <p style="font-size:16px; color:#444; margin:0 0 24px 0; line-height:1.5;">
                  Se precisar de ajuda, nossa equipe está sempre à disposição.
                </p>

                <a href="https://systemeducation.com" 
                  style="display:inline-block; padding:12px 20px; background:#7C3AED; color:#ffffff; text-decoration:none; font-size:16px; font-weight:600; border-radius:6px;">
                  Acessar a plataforma
                </a>

              </td>
            </tr>

            <tr>
              <td style="padding-top:32px; text-align:center; color:#888; font-size:12px;">
                © %d System Education. Todos os direitos reservados.
              </td>
            </tr>

          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`, name, time.Now().Year())

	if err := n.client.SendEmail(emailAddr, subject, html); err != nil {
		return fmt.Errorf("failed to send welcome email to %s: %w", emailAddr, err)
	}

	return nil
}
