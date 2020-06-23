using Microsoft.Win32;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Shapes;

namespace Flipped_Win10
{
    /// <summary>
    /// Register.xaml 的交互逻辑
    /// </summary>
    public partial class Register : Window
    {
        public Register()
        {
            InitializeComponent();
        }

        private void Window_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            DragMove();
        }

        private void CloseButton_Click(object sender, RoutedEventArgs e) 
        {
            Close();
        }

        private void UpLoadImageBtn_Click(object sender, RoutedEventArgs e)
        {
            OpenFileDialog ofd = new OpenFileDialog
            {
                InitialDirectory = @"C:\",
                Filter = "(*.jpg,*.png,*.jpeg,*.bmp,*.gif,*.ico)|*.jgp;*.png;*.jpeg;*.bmp;*.gif;*.ico|All files(*.*)|*.*"
            };
            if (ofd.ShowDialog() == true)
            {
                avatar.Source = new BitmapImage(new Uri(ofd.FileName));
            }
            else
            {
                MessageBox.Show("未选择头像");
            }
        }

        private void SnackbarMessage_ActionClick(object sender, RoutedEventArgs e)
        {
            snackBar.IsActive = false;
        }

        private void RegisterBtn_Click(object sender, RoutedEventArgs e)
        {
            string username = userNameBox.Text;
            string pwd = pwdBox.Password;
            string email = emailBox.Text;
            int userType = userTypeBox.SelectedIndex;
            string avatarSource = avatar.Source == null ? "" : avatar.Source.ToString().Substring(8);

            bool isInputUnvalidated = false; // 判断参数表单是否非法
            string errInfo = "Register Reject: ";

            if (username == null || username.Length == 0)
            {
                isInputUnvalidated = true;
                errInfo += "Username is Empty";
            }
            else if (pwd == null || pwd.Length == 0)
            {
                isInputUnvalidated = true;
                errInfo += "PassWord is Empty";
            }
            else if (email == null || email.Length == 0)
            {
                isInputUnvalidated = true;
                errInfo += "Email Address is Empty";
            }
            else if (!CheckEmail(email))
            {
                isInputUnvalidated = true;
                errInfo += "Email Address is incorrect";
            }
            else if (userType < 0)
            {
                isInputUnvalidated = true;
                errInfo += "userType is not chosen";
            }
            else if (avatarSource == null || avatarSource.Length == 0)
            {
                isInputUnvalidated = true;
                errInfo += "user avatar is Empty";
            }

            if (isInputUnvalidated) 
            {
                snackBarMsg.Content = errInfo;
                snackBar.IsActive = true;
                return;
            }

            var parameters = new Dictionary<String, String>
            {
                { "name", username },
                { "password", pwd },
                { "email", email },
                { "user_type", userType.ToString() },
                { "avatarSource", avatarSource }
            };

            AnalysizeReslut(parameters);
        }

        private async void AnalysizeReslut(Dictionary<String, String> keyValues) 
        {
            var res = await NetWork.RegisterAysnc(keyValues);
            string statusCode = res.Item2;
            string answer = res.Item1;
            if (statusCode == "OK")
            {
                MessageBox.Show("Register Successful!");
                this.Close();
            }
            else
            {
                MessageBox.Show("Register Failed!");
            }
        }

        private static bool CheckEmail(string email) 
        {
            Regex re = new Regex(@"^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$");
            return re.IsMatch(email);
        }
    }
}
