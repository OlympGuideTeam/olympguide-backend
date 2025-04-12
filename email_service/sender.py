import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

import constants
from main import logger


def send_code(to_email, code):
    subject = "OlympGuide - одноразовый код"
    html_body = f"""
    <html>
      <head>
        <style>
          body {{
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
          }}
          .container {{
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
          }}
          .header {{
            color: #2c3e50;
            text-align: center;
            margin-bottom: 25px;
          }}
          .code-container {{
            text-align: center;
            margin: 20px 0;
          }}
          .code {{
            display: inline-block;
            background-color: #f0f0f0;
            border: 1px dashed #ccc;
            padding: 10px 20px;
            font-family: 'Courier New', monospace;
            font-size: 24px;
            letter-spacing: 2px;
            color: #e74c3c;
            border-radius: 5px;
          }}
          .footer {{
            margin-top: 30px;
            font-size: 12px;
            color: #7f8c8d;
            text-align: center;
          }}
          .text-center {{
            text-align: center;
          }}
        </style>
      </head>
      <body>
        <div class="container">
          <div class="header">
            <h2>Ваш одноразовый код подтверждения</h2>
          </div>
          <p class="text-center">Для завершения регистрации введите следующий код:</p>
          <div class="code-container">
            <div class="code">{code}</div>
          </div>
          <p class="text-center">Этот код действителен в течение ограниченного времени.</p>
          <div class="footer">
            <p>OlympGuide &copy; 2025. Все права защищены.</p>
          </div>
        </div>
      </body>
    </html>
    """
    send_mail(html_body, subject, to_email)


def send_password(to_email, password):
    subject = "OlympGuide - сгенерированный пароль"
    html_body = f"""
    <html>
      <head>
        <style>
          body {{
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
          }}
          .container {{
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
          }}
          .header {{
            color: #2c3e50;
            text-align: center;
            margin-bottom: 25px;
          }}
          .password-container {{
            text-align: center;
            margin: 20px 0;
          }}
          .password {{
            display: inline-block;
            background-color: #f8f8f8;
            border: 1px solid #e1e1e1;
            padding: 10px 20px;
            font-family: 'Courier New', monospace;
            font-size: 18px;
            color: #27ae60;
            border-radius: 5px;
            word-break: break-all;
            max-width: 100%;
            overflow-wrap: break-word;
          }}
          .footer {{
            margin-top: 30px;
            font-size: 12px;
            color: #7f8c8d;
            text-align: center;
          }}
          .warning {{
            color: #e74c3c;
            font-weight: bold;
            text-align: center;
          }}
          .text-center {{
            text-align: center;
          }}
        </style>
      </head>
      <body>
        <div class="container">
          <div class="header">
            <h2>Ваш сгенерированный пароль</h2>
          </div>
          <p class="text-center">Вы можете использовать этот пароль для входа по почте:</p>
          <div class="password-container">
            <div class="password">{password}</div>
          </div>
          <p class="warning">Сохраните этот пароль в надежном месте. Для безопасности рекомендуется сменить его после первого входа.</p>
          <div class="footer">
            <p>OlympGuide &copy; 2025. Все права защищены.</p>
          </div>
        </div>
      </body>
    </html>
    """
    send_mail(html_body, subject, to_email)


def send_mail(html_body, subject, to_email):
    msg = MIMEMultipart('alternative')
    msg["Subject"] = subject
    msg["From"] = constants.SMTP_USERNAME
    msg["To"] = to_email

    html_part = MIMEText(html_body, 'html', 'utf-8')
    msg.attach(html_part)

    try:
        with smtplib.SMTP(constants.SMTP_SERVER, constants.SMTP_PORT) as server:
            server.starttls()
            server.login(constants.SMTP_USERNAME, constants.SMTP_PASSWORD)
            server.sendmail(constants.SMTP_USERNAME, to_email, msg.as_string())
        logger.info(f"Email sent to {to_email}")
    except Exception as e:
        logger.error(f"Failed to send email: {e}")