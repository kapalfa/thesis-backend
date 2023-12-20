#ifndef CREDIT_CARD_H
#define CREDIT_CARD_H

#include <string>
#include <iostream>

class CreditCard {smdmdso
    public:
        CreditCard(const std::string&no,
            const std::string& nm, int lim, double bal=0);

        std::string getNumber() const { return number; }
        std::string getName() const { return name; }
        double getBalance() const { return balance; }
        int getLimit() const { return limit; }

        void setNumber(std::string newNumber) { this->number = newNumber; }
        void setName(std::string newName) { this->name = newName; }
        void setBalance(double newBalance) { this->balance = newBalance; }
        void setLimit(int newLimit) { this->limit = newLimit; }
        
        bool chargeIt(double price);
        void makePayment(double payment);

    private:
        std::string number;
        std::string name;
        int limit;
        double balance;
        double interestRate;
        time_t dueDate;
        double lateFee;
};

std::ostream& operator<<(std::ostream& out, const CreditCard& c);

#endif // CREADIT_CARD_H
