#include <cstdlib>
#include <iostream>
#include <string>
using namespace std;
enum MealType { NO_PREF, REGULAR, LOW_FAT, VEGETARIAN };

class Passenger {
    public:
        Passenger(); //default constructor
        Passenger(const string& nm, MealType mp, const string& ffn = "NONE");
        Passenger(const Passenger& pass); //copy constructor
        bool isFrequentFlyer() const;
        void makeFrequentFlyer(const string& newFreqFlyerNo);

    private:
        string name;
        MealType mealPref;
        bool isFreqFlyer;
        string freqFlyerNo;
};

